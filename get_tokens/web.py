#!/usr/bin/env python3
from bottle import redirect, request, response, route, run
import html
from rauth import OAuth1Session

from common import get_service, add_auth_params
from webutil import *

HOST = 'localhost'
PORT = 8090


@route('/')
def index():
    """This route is the main page."""
    response.set_header('Content-Type', 'text/html')
    if cookie_has_access_token():
        # We already have an access token for this client, so show them the
        # form where they can input an API path to request.
        return make_html_page(TEST_FORM)
    else:
        # We don't have an access token for this client, so give them a link to
        # begin the authorization process.
        return make_html_page('<a href="/authorize">Authorize</a>');


@route('/authorize')
def authorize():
    """This route is where our client goes to begin the authorization
    process."""
    service = get_service()

    # Get a request token and secret from SmugMug. This token enables us to
    # make an authorization request.
    rt, rts = service.get_request_token(
            params={'oauth_callback': 'http://localhost:8090/callback'})

    # Record the request token and secret in the client's cookie, because we're
    # going to need it later.
    set_cookie({'rt': rt, 'rts': rts})

    # Get the authorization URL, which is where we send the user so they can
    # approve our authorization request.
    auth_url = service.get_authorize_url(rt)

    # Modify the authorization URL to include the access and permissions levels
    # that our application needs.
    auth_url = add_auth_params(auth_url, access='Full', permissions='Modify')

    # Send the client to the authorization URL.
    redirect(auth_url)


@route('/callback')
def callback():
    """This route is where we receive the callback after the user accepts or
    rejects the authorization request."""
    service = get_service()

    # Get the client's cookie, which contains the request token and secret that
    # we saved earlier.
    cookie = request.get_cookie(COOKIE_NAME, secret=SECRET)

    # Use the request token and secret, and the verifier code received by this
    # callback, to sign the request for an access token.
    at, ats = service.get_access_token(
            cookie['rt'],
            cookie['rts'],
            params={'oauth_verifier': request.query['oauth_verifier']})

    # Store the access token and secret in the client's cookie, replacing the
    # request token and secret (which are no longer valid). The access token
    # and secret will be valid forever, unless the client revokes them.
    set_cookie({'at': at, 'ats': ats})

    # Send the client back to the main page where they can make API requests
    # using their access token and secret.
    redirect('/')


@route('/test')
def test():
    """This route is where our client asks us to make a signed API request on
    their behalf."""
    if not cookie_has_access_token():
        # If we don't have an access token, we can't make a signed API request,
        # so send the client back to the main page, where they can follow the
        # link to reauthorize.
        redirect('/')
        return
    service = get_service()

    # Get the client's cookie, which contains the access token and secret that
    # we saved earlier.
    cookie = request.get_cookie(COOKIE_NAME, secret=SECRET)

    # Make a signed request to the API.
    session = OAuth1Session(
            service.consumer_key, service.consumer_secret,
            access_token=cookie['at'], access_token_secret=cookie['ats'])
    response.set_header('Content-Type', 'text/html')
    json = make_api_request(session, request.query['path'])

    # Show the result of the request to the client.
    return make_html_page(TEST_FORM + '<pre>' + html.escape(json) + '</pre>')


if __name__ == '__main__':
    # Do this now in order to provoke validation errors from config.json.
    get_service()
    run(host=HOST, port=PORT)

