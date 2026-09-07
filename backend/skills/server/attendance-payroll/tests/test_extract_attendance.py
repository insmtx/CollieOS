import json
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1] / "scripts"))
from extract_attendance import json_path_for_page_image, prepare_pages
from pdf_prepare import estimate_skew_angle


class ExtractAttendanceTest(unittest.TestCase):
    def test_copies_image_pages(self):
        with tempfile.TemporaryDirectory() as directory:
            source = Path(directory) / "sheet.png"
            source.write_bytes(b"\x89PNG\r\n")
            output = Path(directory) / "pages"
            pages = prepare_pages([source], output)
            self.assertEqual(len(pages), 1)
            self.assertTrue(Path(pages[0]).is_file())
            self.assertEqual(Path(pages[0]).read_bytes(), source.read_bytes())
            payload = json.dumps({"pages": pages})
            self.assertIn("doc-1.png", payload)

    def test_json_path_sits_beside_rendered_png(self):
        png = Path("/tmp/.attendance-pages/doc-1/page-15-part-1.png")
        self.assertEqual(
            json_path_for_page_image(str(png)),
            str(png.with_name("page-15.json")),
        )

    def test_estimates_small_clockwise_scan_tilt(self):
        from PIL import Image, ImageDraw

        sheet = Image.new("RGB", (420, 280), (255, 255, 255))
        draw = ImageDraw.Draw(sheet)
        for top in range(30, 250, 18):
            draw.rectangle((24, top, 396, top + 3), fill=(20, 20, 20))
        tilted = sheet.rotate(4, expand=True, fillcolor=(255, 255, 255))
        self.assertAlmostEqual(estimate_skew_angle(tilted), -4, delta=1.0)


if __name__ == "__main__":
    unittest.main()
