def parse_csv(filepath, delimiter=','):
    """Parse a CSV file and return a list of dicts.

    Args:
        filepath: Path to the CSV file.
        delimiter: Column separator character.

        Returns:
        list: List of row dictionaries.
    """
    rows = []
    with open(filepath) as f:
        headers = f.readline().strip().split(delimiter)
        for line in f:
            values = line.strip().split(delimiter)
            rows.append(dict(zip(headers, values)))
    return rows


def parse_json(filepath):
    """
    Loads and returns data from a JSON file.
    """
    import json
    with open(filepath) as f:
        return json.load(f)


def validate_schema(data, schema):
    # Checks if data matches the expected schema
    for key, expected_type in schema.items():
        if key not in data:
            raise ValueError(f"Missing key: {key}")
        if not isinstance(data[key], expected_type):
            raise TypeError(f"Expected {expected_type} for {key}")
    return True


def merge_dicts(*dicts):
    """Merge multiple dictionaries into one. Later values override earlier ones."""
    result = {}
    for d in dicts:
        result.update(d)
    return result
