# validator

This is a JSON validator package. It allows tags to be added to structure definitions, and those structures are
then passed to a Define() operation which creates a map of the valid structure definitions. Subsequently, JSON
strings can be validated against the structure definitions to report if they are conformant or not.

This is intended to help catch misspelled fields, missing require fields, and invalid field values.
