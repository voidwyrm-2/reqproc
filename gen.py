#! /usr/bin/env python3
import os
from typing import Any, Callable


def px(msg: str, ec: int = 0):
    print(msg)
    exit(ec)


def smother(fn: Callable, *args: Any) -> Any:
    try:
        return fn(*args)
    except Exception:
        pass


APPNAME = "reqproc"
MAKE = "makefile"
OUTFOLDER = "build"

WINARCH = ["amd64", "arm64"]
LINARCH = ["386", *WINARCH]

TARGETS = {
    "wasip1": ["wasm"],
    "darwin": ["arm64"],
    "windows": WINARCH,
    "linux": LINARCH,
    "freebsd": LINARCH,
    "netbsd": LINARCH,
    "openbsd": LINARCH,
}


def create_commands(targets: dict[str, list[str]]) -> list[str]:
    BPATH = f"{OUTFOLDER}/{APPNAME}"
    return [
        "\n".join(
            [
                f"\tGOOS={tos} GOARCH={tarch} go build -o {BPATH} .\n"
                + f"\tzip -r {BPATH}-{tos}_{tarch}.zip {BPATH}\n"
                + f"\trm -rf {BPATH}"
                for tarch in tarches
            ]
        )
        for tos, tarches in targets.items()
    ]


def generate_targets(targets: dict[str, list[str]]) -> str:
    return (
        "main: clean\n"
        + "\n".join(create_commands(targets))
        + "\n\nlocal: clean\n\t"
        + "\n".join(create_commands({"wasip1": ["wasm"], "darwin": ["arm64"]}))
        + f"\n\nclean:\n\trm -rf {OUTFOLDER}\n\tmkdir {OUTFOLDER}"
    )


with open(MAKE, "wt" if os.path.exists(MAKE) else "xt") as f:
    f.write(generate_targets(TARGETS))
