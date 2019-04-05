SmugMug OAuth Example Code
==========================

This repository contains two example programs written in Python 3 which
demonstrate the use of OAuth with SmugMug.

The programs are:

* `web.py`, which demonstrates OAuth as used by a web application
* `console.py`, which demonstrates how to use OAuth in non-web-application
  scenarios

Before running either example, you must supply a SmugMug API Key and Secret in
a file named `config.json`. The expected format of this file is demonstrated by
`example.json`.

The Web Example
---------------

To run the web example, you can just run `./run-web.sh`. This shell
script will install the Python libraries needed by the example and then run it.

The Non-Web Example
-------------------

To run the non-web example, you can just run `./run-console.sh`. This shell
script will install the Python libraries needed by the example and then run it.
