---
status: Draft
depth: Standard
---

# Feature: Password Reset

## Story
As a registered user, I want to reset my forgotten password via email, so that I can regain access to my account without having to contact support.

## Requirements
- User can submit a password reset request by entering their registered email address on a dedicated page
- System sends a reset email containing a unique, time-limited link (link expires after 1 hour)
- Clicking the link navigates the user to a password reset form
- User can enter a new password and a confirmation field; both must match
- New password must be at least 8 characters long
- After a successful reset, the user is redirected to the login page with a success message
- Used or expired reset tokens are invalidated and cannot be reused
- The reset flow does not reveal whether a given email address is registered (to prevent account enumeration)

## Out of Scope
- Social login / OAuth integration
- Admin-initiated password resets on behalf of users
- SMS-based or TOTP-based reset options
- Password strength meters or complexity rules beyond minimum length

## Open Questions
- Should reset request submissions be rate-limited per email address or per IP?
- Which email service provider will be used to send the reset email?
