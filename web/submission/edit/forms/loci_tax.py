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


class LociTaxonomyForm(Form):

    class TaxonomyForm(Form):
        ncbiTaxId = IntegerField(
            "NCBI Taxonomy ID *", validators=[validators.InputRequired()]
        )
        name = StringField("Name")

    class LocusForm(Form):
        accession = StringField(
            "Genome identifier *",
            [validators.InputRequired(), ValidateSingleInput(), validate_genbank],
            widget=TextInputIndicator(),
            description="E.g., AL645882. Only use GenBank accessions, not RefSeq accessions or GI numbers.",
            # render_kw={
            #     "hx-post": "/query_ncbi",
            #     "hx-trigger": "change",
            #     "hx-swap": "innerHTML",
            #     "hx-target": ".subform#taxonomy",
            #     "hx-indicator": "#spinner",
            # },
        )
        draft_genome = BooleanField(
            "This accession is a draft",
            description="Select if this accession is not yet publicised on GenBank",
        )
        location = FormField(
            location_form_factory(),
            description="Start and end coordinates, may be left empty if gene cluster spans entire record.",
        )
        evidence = FieldList(
            FormField(EvidenceForm),
            min_entries=1,
            description="Type of evidence that shows this gene cluster is responsible for the biosynthesis of the designated molecule. Papers highlighting multiple methods can be added under each applicable method using the 'Add additional evidence' button.",
            widget=FieldListAddBtn(
                label="Add additional evidence",
            ),
        )

    loci = FieldList(
        FormField(LocusForm),
        min_entries=1,
        description="Locus or loci where the gene cluster is located",
        widget=FieldListAddBtn(
            label="Add additional locus",
        ),
    )

    taxonomy = FormField(TaxonomyForm)
