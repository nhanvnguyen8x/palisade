from __future__ import annotations

import re
import uuid
from typing import Any


def _normalize_text(text: str) -> str:
    # Collapse whitespace for more stable chunk boundaries.
    text = re.sub(r"[ \t]+", " ", text)
    text = re.sub(r"\n{3,}", "\n\n", text)
    return text.strip()


def chunk_pages(
    pages: list[str],
    max_chars: int,
) -> list[dict[str, Any]]:
    """
    Dev chunking strategy:
    - Use page-level text.
    - If a page is too large, split by whitespace into sub-chunks of ~max_chars.
    """
    chunks: list[dict[str, Any]] = []
    chunk_index = 0

    for page_number, page_text in enumerate(pages):
        page_text = _normalize_text(page_text)
        if not page_text:
            continue

        if len(page_text) <= max_chars:
            chunks.append(
                {
                    "chunk_index": chunk_index,
                    "page_number": page_number,
                    "text": page_text,
                    "chunk_id": uuid.uuid4(),
                    "embedding_id": uuid.uuid4(),
                }
            )
            chunk_index += 1
            continue

        # Split into smaller chunks
        words = page_text.split(" ")
        buf: list[str] = []
        buf_len = 0

        def flush() -> None:
            nonlocal chunk_index, buf, buf_len
            if not buf:
                return
            txt = " ".join(buf).strip()
            if txt:
                chunks.append(
                    {
                        "chunk_index": chunk_index,
                        "page_number": page_number,
                        "text": txt,
                        "chunk_id": uuid.uuid4(),
                        "embedding_id": uuid.uuid4(),
                    }
                )
                chunk_index += 1
            buf = []
            buf_len = 0

        for w in words:
            if buf_len + len(w) + 1 > max_chars and buf:
                flush()
            buf.append(w)
            buf_len += len(w) + 1
        flush()

    return chunks

