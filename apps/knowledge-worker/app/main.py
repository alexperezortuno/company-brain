from __future__ import annotations

from pathlib import Path
from typing import Any

import yaml
from fastapi import FastAPI
from pydantic import BaseModel
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="", case_sensitive=False)

    brain_config_path: str = "/app/config/brain.yaml"


class HealthResponse(BaseModel):
    status: str
    service: str
    instance: str
    modules: dict[str, Any]


settings = Settings()


def load_brain_config() -> dict[str, Any]:
    path = Path(settings.brain_config_path)
    with path.open("r", encoding="utf-8") as stream:
        data = yaml.safe_load(stream) or {}
    instance_name = data.get("instance", {}).get("name")
    if not instance_name:
        raise RuntimeError("instance.name is required")
    return data


brain_config = load_brain_config()
app = FastAPI(title="Company Brain Knowledge Worker", version="0.1.0")


@app.get("/live")
def live() -> dict[str, str]:
    return {"status": "alive"}


@app.get("/health", response_model=HealthResponse)
def health() -> HealthResponse:
    return HealthResponse(
        status="healthy",
        service="knowledge-worker",
        instance=brain_config["instance"]["name"],
        modules=brain_config.get("modules", {}),
    )
