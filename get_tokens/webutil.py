"""This module contains utilities used by the SmugMug OAuth Demo which are not
relevant to understanding the OAuth workflow."""
from bottle import request, response
from urllib.parse import urlparse, parse_qs
import uuid

from common import API_ORIGIN

COOKIE_NAME = 'c'

# This secret is used for signing cookies.
# If you want the cookies to be valid from one run of this program to the next,
# use a static string here (but make sure to choose it randomly).
SECRET = str(uuid.uuid4())

def get_cookie():
    """Get the client's cookie, which is where we store the client's request
    token or access token.

    Note that the cookie is not sent to SmugMug, and using OAuth does not
    require that you use cookies.  It is simply the mechanism by which this
    particular example application keeps track of the state of each of its
    clients.
    """
    return request.get_cookie(COOKIE_NAME, secret=SECRET)


def cookie_has_access_token():
    """Does the client's cookie contain an access token?"""
    c = get_cookie()
    return type(c) is dict and 'at' in c and 'ats' in c


def set_cookie(obj):
    """Set the contents of the client's cookie."""
    response.set_cookie(
            COOKIE_NAME, obj, secret=SECRET, httponly=True, path='/')


def make_html_page(content):
    return '<html><head><title>SmugMug OAuth Demo</title></head><body>' \
            + content + '</body></html>'


TEST_FORM = '''
    <form action="/test" method="GET">
        <label for="path">API URI Path:</label>
        <input type="text" name="path" id="path" value="/api/v2!authuser">
        <input type="submit" value="GET">
    </form>
'''


def make_api_request(session, input_path):
    parsed = urlparse(input_path)
    params = parse_qs(parsed.query)
    params['_pretty'] = ''
    return session.get(
            API_ORIGIN + parsed.path,
            params=params,
            headers={'Accept': 'application/json'}).text

