from flask import current_app, session
from markupsafe import Markup
import requests
from wtforms import (
    DecimalField,
    Form,
    HiddenField,
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
from wtforms.widgets import HiddenInput
from submission.utils.custom_fields import ReferenceField, TagListField
from submission.utils.custom_forms import location_form_factory, LociEvidenceForm
from submission.utils.custom_widgets import (
    FieldListAddBtn,
    TextInputIndicator,
    SubmitIndicator,
    SelectDefault,
    ProductInputSearch,
)
from submission.utils.custom_validators import (
    ValidateCitations,
    ValidateSingleInput,
    validate_genbank,
    validate_loci,
)


class BioActivityAssayMeasurementSubForm(Form):
    db_id = IntegerField(widget=HiddenInput())
    concentration = DecimalField(places=None)
    unit = StringField()
    error = DecimalField(places=None)
    replicates = DecimalField(places=None)


class BioActivityAssayTestSystemSubForm(Form):
    db_id = IntegerField(widget=HiddenInput())
    cell_line = StringField()
    organism = IntegerField()
    strain = StringField()


class BioActivityAssaySubForm(Form):
    db_id = IntegerField(widget=HiddenInput())
    db_bio_activity_id = IntegerField(widget=HiddenInput())
    type = SelectField(
        choices=[
            "IC50",
            "EC50",
            "GI50",
            "CC50",
            "MIC",
            "MBC",
            "Ki",
            "Kd",
            "Km",
            "% inhibition at concentration",
            "Zone of inhibition",
            "Fold change",
            "Other",
        ]
    )
    db_measurement_id = IntegerField(widget=HiddenInput())
    measurement = FormField(BioActivityAssayMeasurementSubForm)
    target = SelectField(
        choices=[
            "Molecular target",
            "Enzyme",
            "Receptor",
            "Whole organism",
            "Cell line",
            "Primary cells",
            "Tissue",
            "Organoid",
            "In vivo model",
        ]
    )
    details = StringField()
    db_test_system_id = IntegerField(widget=HiddenInput())
    test_system = FormField(BioActivityAssayTestSystemSubForm)
    references = ReferenceField(
        label="Citation(s) *",
        description=Markup(
            "Accepted formats are (in order of preference):<br>"
            "'doi:10.1016/j.chembiol.2020.11.009', 'pubmed:33321099', 'patent:US7070980B2', "
            "'url:https://example.com'.<br>If no publication "
            "is available <u>yet</u>, please use 'doi:pending'."
        ),
        validators=[validators.InputRequired(), ValidateCitations()],
    )


class BioactivitySubForm(Form):
    db_id = IntegerField(widget=HiddenInput())
    compound_id = IntegerField(widget=HiddenInput())
    name = SelectField(
        choices=[
            "Cytotoxicity",
            "Antiproliferative",
            "Antibiotic (antibacterial)",
            "Antifungal",
            "Antiviral",
            "Antiparasitic",
            "Anti-inflammatory",
            "Enzyme inhibition",
            "Receptor agonist",
            "Receptor antagonist",
            "Ion channel modulator",
            "Neuroprotective",
            "Antioxidant",
            "Immunomodulatory",
            "Other",
        ]
    )
    details = StringField()
    observed = BooleanField()
    assays = FieldList(
        FormField(BioActivityAssaySubForm),
        widget=FieldListAddBtn(
            label="Add bioactivity assay",
        ),
    )
    references = ReferenceField(
        label="Citation(s) *",
        description=Markup(
            "Accepted formats are (in order of preference):<br>"
            "'doi:10.1016/j.chembiol.2020.11.009', 'pubmed:33321099', 'patent:US7070980B2', "
            "'url:https://example.com'.<br>If no publication "
            "is available <u>yet</u>, please use 'doi:pending'."
        ),
        validators=[validators.InputRequired(), ValidateCitations()],
    )


class CompoundEvidence(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0, validators=None)
    db_compound_id = IntegerField(widget=HiddenInput(), default=0)
    method = StringField()
    references = ReferenceField(
        label="Citation(s) *",
        description=Markup(
            "Accepted formats are (in order of preference):<br>"
            "'doi:10.1016/j.chembiol.2020.11.009', 'pubmed:33321099', 'patent:US7070980B2', "
            "'url:https://example.com'.<br>If no publication "
            "is available <u>yet</u>, please use 'doi:pending'."
        ),
        validators=[validators.InputRequired(), ValidateCitations()],
    )


class CompoundsSubForm(Form):
    db_id = IntegerField(widget=HiddenInput(), default=0)
    name = StringField()
    evidence = FieldList(
        FormField(CompoundEvidence),
        widget=FieldListAddBtn(
            label="Add evidence",
        ),
        min_entries=1,
    )
    bioactivities = FieldList(
        FormField(BioactivitySubForm),
        widget=FieldListAddBtn(
            label="Add bioactivity",
        ),
    )
    structure = StringField()
    databaseIds = FieldList(
        StringField(),
        label="Database IDs",
        widget=FieldListAddBtn(
            label="Add database",
        ),
    )
    mass = DecimalField("Exact Mass", default=0, places=None)
    formula = StringField("Molecular Formula")


class CompoundsForm(Form):
    compounds = FieldList(
        FormField(CompoundsSubForm),
        widget=FieldListAddBtn(
            label="Add compound",
        ),
    )
