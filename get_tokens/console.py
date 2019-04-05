#!/usr/bin/env python3
from rauth import OAuth1Session
import sys

from common import API_ORIGIN, get_service, add_auth_params


def main():
    """This example interacts with its user through the console, but it is
    similar in principle to the way any non-web-based application can obtain an
    OAuth authorization from a user."""
    service = get_service()

    # First, we need a request token and secret, which SmugMug will give us.
    # We are specifying "oob" (out-of-band) as the callback because we don't
    # have a website for SmugMug to call back to.
    rt, rts = service.get_request_token(params={'oauth_callback': 'oob'})

    # Second, we need to give the user the web URL where they can authorize our
    # application.
    auth_url = add_auth_params(
            service.get_authorize_url(rt), access='Full', permissions='Modify')
    print('Go to %s in a web browser.' % auth_url)

    # Once the user has authorized our application, they will be given a
    # six-digit verifier code. Our third step is to ask the user to enter that
    # code:
    sys.stdout.write('Enter the six-digit code: ')
    sys.stdout.flush()
    verifier = sys.stdin.readline().strip()

    # Finally, we can use the verifier code, along with the request token and
    # secret, to sign a request for an access token.
    at, ats = service.get_access_token(rt, rts, params={'oauth_verifier': verifier})

    # The access token we have received is valid forever, unless the user
    # revokes it.  Let's make one example API request to show that the access
    # token works.
    print('\n\nAccess token: %s' % at)
    print('Access token secret: %s\n\n' % ats)
    session = OAuth1Session(
            service.consumer_key,
            service.consumer_secret,
            access_token=at,
            access_token_secret=ats)
    print('Example response:\n')
    print(session.get(
        API_ORIGIN + '/api/v2!authuser',
        headers={'Accept': 'application/json'}).text)


if __name__ == '__main__':
    main()

