from datetime import datetime, timezone


def iso_format(dt: datetime) -> str:
    dt = dt.astimezone(timezone.utc)
    ms = f"{dt.microsecond // 1000:03d}"
    return dt.strftime(f"%Y-%m-%dT%H:%M:%S.{ms}Z")


def iso_now() -> str:
    return iso_format(datetime.now(timezone.utc))
