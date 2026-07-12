import os


class Config:
    # Claim/poll loop
    poll_interval_seconds: float
    batch_size: int

    # Postgres
    database_url: str

    # MinIO
    minio_endpoint: str
    minio_access_key: str
    minio_secret_key: str
    minio_bucket: str
    minio_use_ssl: bool

    # Indexing
    embedding_dim: int
    chunk_max_chars: int

    # Optional Milvus
    milvus_enabled: bool
    milvus_host: str
    milvus_port: int
    milvus_collection: str

    def __init__(self) -> None:
        self.poll_interval_seconds = float(os.getenv("POLL_INTERVAL_SECONDS", "2"))
        self.batch_size = int(os.getenv("INGESTION_BATCH_SIZE", "5"))

        # Reuse Go env vars for consistency.
        # If DATABASE_URL is not provided, build it from DB_* vars.
        self.database_url = os.getenv("DATABASE_URL", "").strip()
        if not self.database_url:
            db_host = os.getenv("DB_HOST", "localhost")
            db_port = int(os.getenv("DB_PORT", "5432"))
            db_user = os.getenv("DB_USER", "palisade")
            db_password = os.getenv("DB_PASSWORD", "palisade")
            db_name = os.getenv("DB_NAME", "palisade")
            sslmode = os.getenv("DB_SSLMODE", "disable")

            # pgx uses sslmode strings; for psycopg it can be passed via querystring.
            # sslmode=disable => sslmode=disable (no TLS).
            self.database_url = (
                f"postgresql://{db_user}:{db_password}@{db_host}:{db_port}/{db_name}"
                f"?sslmode={sslmode}"
            )

        self.minio_endpoint = os.getenv("MINIO_ENDPOINT", "localhost:9000")
        self.minio_access_key = os.getenv("MINIO_ACCESS_KEY", "admin")
        self.minio_secret_key = os.getenv("MINIO_SECRET_KEY", "password123")
        self.minio_bucket = os.getenv("MINIO_BUCKET", "palisade")
        self.minio_use_ssl = os.getenv("MINIO_USE_SSL", "false").lower() == "true"

        self.embedding_dim = int(os.getenv("EMBEDDING_DIM", "384"))
        self.chunk_max_chars = int(os.getenv("CHUNK_MAX_CHARS", "2000"))

        self.milvus_enabled = os.getenv("MILVUS_ENABLED", "false").lower() == "true"
        self.milvus_host = os.getenv("MILVUS_HOST", "localhost")
        self.milvus_port = int(os.getenv("MILVUS_PORT", "19530"))
        self.milvus_collection = os.getenv("MILVUS_COLLECTION", "knowledge_chunks")

