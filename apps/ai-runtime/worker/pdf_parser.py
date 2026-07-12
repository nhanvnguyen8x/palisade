from __future__ import annotations

import io
from typing import List

from pypdf import PdfReader


def extract_text_by_page(pdf_bytes: bytes) -> List[str]:
    reader = PdfReader(io.BytesIO(pdf_bytes))

    pages: List[str] = []
    for page in reader.pages:
        text = page.extract_text() or ""
        pages.append(text)

    return pages

