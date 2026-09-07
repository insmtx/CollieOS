"""Prepare bounded, OCR-friendly page images from a PDF."""
import os
import sys
from io import BytesIO

import fitz
from PIL import Image

MAX_IMAGE_SIDE = 1800
MAX_IMAGE_BYTES = 3 * 1024 * 1024
MAX_DESKEW_DEGREES = 8.0
DESKEW_STEP = 0.5
DESKEW_MIN_DEGREES = 0.4


def non_white_bounds(pix):
    """Return the content bounds, preserving a small margin for table lines."""
    samples, width, height, stride = pix.samples, pix.width, pix.height, pix.stride
    rows = []
    for y in range(height):
        row = samples[y * stride:(y + 1) * stride]
        if any(value < 245 for value in row):
            rows.append(y)
    if not rows:
        return 0, height
    return max(0, rows[0] - 12), min(height, rows[-1] + 13)


def _as_rgb(image):
    if image.mode == "RGB":
        return image
    if image.mode == "RGBA":
        background = Image.new("RGB", image.size, (255, 255, 255))
        background.paste(image, mask=image.split()[-1])
        return background
    return image.convert("RGB")


def _projection_score(gray):
    width, height = gray.size
    pixels = gray.tobytes()
    row_ink = []
    for y in range(height):
        start = y * width
        row = pixels[start:start + width]
        row_ink.append(sum(1 for value in row[::2] if value < 200))
    if not row_ink:
        return 0.0
    mean = sum(row_ink) / len(row_ink)
    return sum((value - mean) ** 2 for value in row_ink) / len(row_ink)


def estimate_skew_angle(image, max_angle=MAX_DESKEW_DEGREES, step=DESKEW_STEP):
    """Find a small rotation that makes table rows more horizontal."""
    gray = _as_rgb(image).convert("L")
    width, height = gray.size
    longest = max(width, height)
    if longest > 720:
        scale = 720 / longest
        gray = gray.resize(
            (max(1, int(width * scale)), max(1, int(height * scale))),
            Image.Resampling.BILINEAR,
        )
    best_angle = 0.0
    best_score = -1.0
    angle = -max_angle
    while angle <= max_angle + 1e-9:
        rotated = gray.rotate(angle, resample=Image.Resampling.BILINEAR, fillcolor=255)
        score = _projection_score(rotated)
        if score > best_score:
            best_score = score
            best_angle = angle
        angle += step
    return best_angle


def deskew_image(image):
    """Rotate a mildly tilted scan so name rows stay aligned with 实际出勤."""
    rgb = _as_rgb(image)
    angle = estimate_skew_angle(rgb)
    if abs(angle) < DESKEW_MIN_DEGREES:
        return rgb
    return rgb.rotate(
        angle,
        resample=Image.Resampling.BICUBIC,
        expand=True,
        fillcolor=(255, 255, 255),
    )


def render_pdf(source, target, limit):
    """Render each original PDF page to one bounded image.

    The page remains intact: rows and date columns are never split. The size
    limit prevents OpenCode from entering its large-image normalization path
    when the model reads a generated page image. Mild camera/scan tilt is
    corrected before the model reads 实际出勤, which sits on the far right.
    """
    document = fitz.open(source)
    if len(document) > limit:
        raise RuntimeError("PDF exceeds page limit")
    paths = []
    for page_number, page in enumerate(document, 1):
        full = page.get_pixmap(dpi=120, alpha=False)
        start, end = non_white_bounds(full)
        image = Image.frombytes("RGB", (full.width, full.height), full.samples)
        image = deskew_image(image.crop((0, start, image.width, end)))
        image.thumbnail((MAX_IMAGE_SIDE, MAX_IMAGE_SIDE), Image.Resampling.LANCZOS)

        png_buffer = BytesIO()
        image.save(png_buffer, format="PNG", optimize=True)
        suffix = ".png"
        payload = png_buffer.getvalue()
        if len(payload) > MAX_IMAGE_BYTES:
            # Keep text/table edges readable while avoiding OpenCode's
            # base64/image.normalize threshold. PNG remains the default.
            for quality in (85, 75, 65):
                jpeg_buffer = BytesIO()
                image.save(jpeg_buffer, format="JPEG", quality=quality, optimize=True)
                payload = jpeg_buffer.getvalue()
                suffix = ".jpg"
                if len(payload) <= MAX_IMAGE_BYTES:
                    break

        path = os.path.join(target, f"page-{page_number}-part-1{suffix}")
        with open(path, "wb") as output:
            output.write(payload)
        paths.append(path)
    return paths


if __name__ == "__main__":
    source, target, limit = sys.argv[1], sys.argv[2], int(sys.argv[3])
    render_pdf(source, target, limit)
