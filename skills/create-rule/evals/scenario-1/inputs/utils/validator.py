def validate_email(email):
    """Returns True if the email is valid, False otherwise."""
    import re
    pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    return bool(re.match(pattern, email))


def validate_phone(phone):
    """
    Validate a phone number string.

    Parameters
    ----------
    phone : str
        The phone number to validate.

    Returns
    -------
    bool
        True if valid.
    """
    import re
    return bool(re.match(r'^\+?[\d\s\-]{7,15}$', phone))


def sanitize_input(value, max_length=255):
    """Sanitize user input by stripping whitespace and truncating."""
    if not isinstance(value, str):
        raise TypeError("Input must be a string")
    return value.strip()[:max_length]


def validate_date(date_str):
    # Accepts dates in YYYY-MM-DD format
    import re
    return bool(re.match(r'^\d{4}-\d{2}-\d{2}$', date_str))


def is_positive_integer(value):
    """
    Check whether value is a positive integer.

    :param value: The value to check.
    :type value: any
    :return: True if value is a positive integer, False otherwise.
    :rtype: bool
    """
    return isinstance(value, int) and value > 0
