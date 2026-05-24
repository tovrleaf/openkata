def format_currency(amount, currency='USD'):
    """Format a numeric amount as a currency string."""
    symbols = {'USD': '$', 'EUR': '€', 'GBP': '£'}
    symbol = symbols.get(currency, currency)
    return f"{symbol}{amount:,.2f}"


def truncate_text(text, max_chars, suffix='...'):
    # Truncates text to max_chars, appending suffix if truncated
    if len(text) <= max_chars:
        return text
    return text[:max_chars - len(suffix)] + suffix


def format_list(items, separator=', ', last_separator=' and '):
    """
    Formats a list of items as a human-readable string.
    For example, ['a', 'b', 'c'] becomes 'a, b and c'.
    The separator is used between all items except the last two.
    The last_separator is used between the final two items.
    """
    if not items:
        return ''
    if len(items) == 1:
        return items[0]
    return separator.join(items[:-1]) + last_separator + items[-1]


def pad_string(s, width, fill=' ', align='left'):
    """Pad a string to the given width using fill character."""
    if align == 'left':
        return s.ljust(width, fill)
    elif align == 'right':
        return s.rjust(width, fill)
    return s.center(width, fill)
