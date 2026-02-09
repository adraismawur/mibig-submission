from dataclasses import dataclass

from flask import current_app, session
import requests
from submission.edit.forms.bio_synth import BioSynthForm
from submission.edit.forms.compounds import CompoundsForm
from submission.edit.forms.finalize import FinalizeForm
from submission.edit.forms.gene_annotation import GeneAnnotationForm
from submission.edit.forms.loci_tax import LociTaxonomyForm


@dataclass
class WizardPage:
    id: str
    description: str
    form: type = None
    data_get_endpoint: str = "/entry/<bgc_id>"
    data_set_endpoint: str = "/entry/<bgc_id>"
    template: str = "wizard/main.html"

    def create_form(self, request_form, entry):
        return self.form(request_form, data=entry)

    def get_data(self, bgc_id):
        replaced_api_endpoint = self.data_get_endpoint.replace("<bgc_id>", bgc_id)

        response = requests.get(
            f"{current_app.config['API_BASE']}" + replaced_api_endpoint,
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            data = response.json()

            return data
        return None
    
    def post_data(self, bgc_id, data: dict[str, any]):
        replaced_api_endpoint = self.data_set_endpoint.replace("<bgc_id>", bgc_id)

        response = requests.post(
            f"{current_app.config['API_BASE']}" + replaced_api_endpoint,
            headers={"Authorization": f"Bearer {session['token']}"},
            data=data
        )

        return response.status_code == 200
            


wizard_pages = [
    WizardPage(
        "locitax",
        "basic information",
        LociTaxonomyForm,
        data_set_endpoint="/entr/<bgc_id>/locitax",
    ),
    WizardPage(
        "biosynth",
        "biosynthetic information",
        BioSynthForm,
        data_get_endpoint="/entry/<bgc_id>/biosynth",
        template="wizard/biosynth.html",
    ),
    WizardPage("compounds", "compound information", CompoundsForm),
    WizardPage("gene_annotation", "additional gene annotation", GeneAnnotationForm),
    WizardPage("finalize", "final details", FinalizeForm),
]

wizard_page_index = {page.id: i for i, page in enumerate(wizard_pages)}


def get_default_wizard_page():
    return wizard_pages[0].id


def get_wizard_page(page_id: str) -> WizardPage:
    if page_id not in wizard_page_index:
        return None

    form_idx = wizard_page_index[page_id]

    if form_idx < 0 or form_idx > len(wizard_pages) - 1:
        return None

    return wizard_pages[form_idx]


def get_next_page(page_id: str) -> type:
    next_idx = wizard_page_index[page_id] + 1

    if next_idx > len(wizard_pages) - 1:
        return None

    return wizard_pages[next_idx]


def get_prev_page(page_id: str) -> tuple[int, type]:
    prev_idx = wizard_page_index[page_id] - 1

    if prev_idx < 0:
        return None

    return wizard_pages[prev_idx]
