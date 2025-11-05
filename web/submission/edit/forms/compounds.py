from flask import current_app, session
import requests
from wtforms import (
    Form,
    StringField,
    IntegerField,
    FieldList,
    SelectField,
    SelectMultipleField,
    FormField,
    BooleanField,
    ValidationError,
    validators,
    SubmitField,
)
from submission.utils.custom_fields import TagListField
from submission.utils.custom_forms import location_form_factory, EvidenceForm
from submission.utils.custom_widgets import (
    FieldListAddBtn,
    TextInputIndicator,
    SubmitIndicator,
    SelectDefault,
    ProductInputSearch,
)
from submission.utils.custom_validators import (
    ValidateSingleInput,
    validate_genbank,
    validate_loci,
)


class BioactivitySubForm(Form):
    name = StringField()
    observed = BooleanField()
    references = FieldList(
        StringField(),
        widget=FieldListAddBtn(
            label="Add reference",
        ),
    )


class CompoundsSubForm(Form):
    name = StringField()
    bioactivities = FieldList(
        FormField(BioactivitySubForm),
        widget=FieldListAddBtn(
            label="Add bioactivity",
        ),
    )
    structure = StringField()
    databaseIds = FieldList(
        StringField(),
        widget=FieldListAddBtn(
            label="Add database",
        ),
    )


class CompoundsForm(Form):
    compounds = FieldList(
        FormField(CompoundsSubForm),
        widget=FieldListAddBtn(
            label="Add compound",
        ),
    )
