#!/bin/sh
if [ ! -d ENV ]; then
    pyvenv ENV
fi
. ENV/bin/activate
pip install -r requirements.txt
