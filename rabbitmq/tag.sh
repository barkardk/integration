#!/bin/bash
git describe --tags --dirty --match='v*' 2>/dev/null || echo v0.0.0 | cut -c2-
