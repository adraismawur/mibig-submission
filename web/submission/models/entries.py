from typing import Any, Union

from flask import current_app, session
import requests
from sqlalchemy import select
from sqlalchemy.orm import Mapped, mapped_column, relationship

from submission.extensions import db
from submission.models import Reference


# TODO: expand table
class Entry(db.Model):
    __tablename__ = "entries"
    __table_args__ = {"schema": "edit"}

    id: Mapped[int] = mapped_column(autoincrement=True, primary_key=True)
    identifier: Mapped[str]
    references: Mapped[list["Reference"]] = relationship(
        "Reference",
        secondary="edit.entry_references",
        back_populates="entries",
        lazy="selectin",
    )

    @classmethod
    def create(cls, bgc_id: str):
        """Create a database entry for this BGC

        Args:
            bgc_id (str): BGC identifier
        """
        entry = cls(identifier=bgc_id)
        db.session.add(entry)
        db.session.commit()
        return entry
    
    def mutate(from_accession: str):
        return requests.post(
            f"{current_app.config['API_BASE']}/mutation",
            headers={"Authorization": f"Bearer {session['token']}"},
            json={
                "from_accession": from_accession
            }
        )

    def add_references(self, refs: list["Reference"]):
        """Add references to this entry if they are not already present

        Args:
            refs (list["Reference"]): list of Reference database objects
        """
        for ref in refs:
            if ref not in self.references:
                self.references.append(ref)
        db.session.commit()

    @staticmethod
    def get(bgc_id: str) -> Union["Entry", None]:
        """Get an entry from database based on identifier

        Args:
            bgc_id (str): BGC identifier

        Returns:
            Entry | None: entry database object or none if not exists
        """
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            entry = response.json()

            return entry
        return None

    def get_text(bgc_id: str):
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}?pretty=true",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            return response.text

        return None

    def get_class(bgc_id: str, class_id: int):
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/class/{class_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            biosynthetic_class = response.json()

            biosynthetic_class["class_"] = biosynthetic_class["class"]

            # 'other' overlaps on 'other' module and 'other' class. fix that here
            if biosynthetic_class["class"] == "other":
                biosynthetic_class["class"] = "class_other"

            return biosynthetic_class
        return None

    def get_class_text(bgc_id: str, class_id: int):
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/class/{class_id}?pretty=true",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            return response.text

        return None

    def create_class(bgc_id: str, data: dict[str, any]):

        if "class_" in data:
            data["class"] = data["class_"]
            # 'other' overlaps on 'other' module and 'other' class. fix that here
            if data["class"] == "class_other":
                data["class"] = "other"

        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/class",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data,
        )
        if response.status_code == 200:
            return True

        return False

    def update_class(bgc_id: str, class_id: int, data: dict[str, any]):

        if "class_" in data:
            data["class"] = data["class_"]
            # 'other' overlaps on 'other' module and 'other' class. fix that here
            if data["class"] == "class_other":
                data["class"] = "other"

        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/class/{class_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data,
        )
        if response.status_code == 200:
            return None

        return response.json()["error"]

    def delete_class(bgc_id: str, class_id: int):
        response = requests.delete(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/class/{class_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            return True

        return False

    def get_module(bgc_id: str, name: int):
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/module/{name}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            data = response.json()

            if data["type"] == "other":
                data["type"] = "module_other"

            data["type"] = data["type"].replace("-", "_")

            return data
        return None

    def move_module(bgc_id, id_from, id_to):
        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/module_reorder",
            headers={"Authorization": f"Bearer {session['token']}"},
            json={"id_from": int(id_from), "id_to": int(id_to)},
        )
        if response.status_code == 200:
            return None

        return response.json()["error"]

    def get_module_text(bgc_id: str, module_id: str):
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/module/{module_id}?pretty=true",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            return response.text

        return None

    def create_module(bgc_id: str, data: dict[str, any]):
        if data["type"] == "module_other":
            data["type"] = "other"

        data["type"] = data["type"].replace("_", "-")

        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/module",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data,
        )
        if response.status_code == 200:
            return True

        return False

    def update_module(bgc_id: str, module_id: str, data: dict[str, any]):
        if data["type"] == "module_other":
            data["type"] = "other"

        data["type"] = data["type"].replace("_", "-")

        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/module/{module_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data,
        )
        if response.status_code == 200:
            return (True, None)

        return (False, response.json()['error'])

    def delete_module(bgc_id: str, module_id: str):
        response = requests.delete(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/module/{module_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            return True

        return False
    

    def get_modification_domain_list(bgc_id: str, module_id: int):
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/modification_domain/list/{module_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            data = response.json()

            for domain in data:
                domain["domain_type"] = domain["type"]

            return (data, None)
        
        return (None, response.json()['error'])
    

    def get_modification_domain(bgc_id: str, modification_domain_id: int, pretty=False):
        url = f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/modification_domain/{modification_domain_id}"

        if pretty:
            url = url + "?pretty=true"

        response = requests.get(
            url,
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            if pretty:
                return response.text, None
            else:
                data = response.json()
                data["domain_type"] = data["type"]
                return (response.json(), None)
        
        return (None, response.json()['error'])
    
    def create_modification_domain(bgc_id: str, module_id: int, data):
        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/modification_domain/add/{module_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data
        )

        if response.status_code == 200:
            return (True, None)
        
        return (False, response.json()['error'])

    def update_modification_domain(bgc_id: str, modification_domain_id: int, data):
        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/modification_domain/{modification_domain_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data
        )

        if response.status_code == 200:
            data = response.json()
            return (data, None)
        
        return (None, response.json()['error'])

    def remove_modification_domain(bgc_id: str, modification_domain_id: int):
        response = requests.delete(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/modification_domain/{modification_domain_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return (True, None)
        
        return (False, response.json()['error'])
        



    def get_path(bgc_id: str, path_id: int):
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/path/{path_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            data = response.json()

            return data
        return None

    def get_path_text(bgc_id: str, path_id: id):
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/path/{path_id}?pretty=true",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            return response.text

        return None

    def create_path(bgc_id: str, biosynth_id: int, data: dict[str, any]):

        data["db_biosynth_id"] = biosynth_id

        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/path",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data,
        )
        if response.status_code == 200:
            return True

        return False

    def update_path(bgc_id: str, path_id: str, data: dict[str, any]):
        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/path/{path_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data,
        )
        if response.status_code == 200:
            return True

        return False

    def delete_path(bgc_id: str, path_id: str):
        response = requests.delete(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/biosynth/path/{path_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            return True

        return False

    def get_compound(bgc_id: str, id: int = None):
        request_url = f"{current_app.config['API_BASE']}/entry/{bgc_id}/compounds"

        if id:
            request_url = request_url + "?id=" + str(id)

        response = requests.get(
            request_url,
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return response.json()

        return None

    def update_compound(bgc_id: str, compound_data: dict[any]):
        compound_data["mass"] = float(compound_data["mass"])

        request_url = f"{current_app.config['API_BASE']}/entry/{bgc_id}/compounds/{compound_data['db_id']}"

        response = requests.post(
            request_url,
            headers={"Authorization": f"Bearer {session['token']}"},
            json=compound_data,
        )

        return response

    def create_compound(bgc_id: str, compound_data: dict[any]):
        compound_data["mass"] = float(compound_data["mass"])

        request_url = f"{current_app.config['API_BASE']}/entry/{bgc_id}/compounds"

        response = requests.post(
            request_url,
            headers={"Authorization": f"Bearer {session['token']}"},
            json=compound_data,
        )

        return response.json()

    def get_compound_text(bgc_id: str, compound_id: int):
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/compounds?id={compound_id}&pretty=true",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return response.text

        return None

    def delete_compound(bgc_id: str, compound_id: int):
        response = requests.delete(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/compounds/{compound_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )
        if response.status_code == 200:
            return True

        return False

    def get_gene_addition(bgc_id: str, addition_id: int, pretty=False):
        pretty = str(pretty).lower()
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/to_add/{addition_id}?pretty={pretty}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return response.json(), None

        return None, response.json()["error"]

    def get_gene_deletion(bgc_id: str, deletion_id: int, pretty=False):
        pretty = str(pretty).lower()
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/to_delete/{deletion_id}?pretty={pretty}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return response.json(), None

        return None, response.json()["error"]

    def get_gene_annotation(bgc_id: str, annotation_id: int, pretty=False):
        pretty = str(pretty).lower()
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/annotation/{annotation_id}?pretty={pretty}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return response.json(), None

        return None, response.json()["error"]

    def get_gene_addition_text(bgc_id: str, addition_id: int, pretty=False):
        pretty = str(pretty).lower()
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/to_add/{addition_id}?pretty={pretty}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return response.text, None

        return None, response.json()["error"]

    def get_gene_deletion_text(bgc_id: str, deletion_id: int, pretty=False):
        pretty = str(pretty).lower()
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/to_delete/{deletion_id}?pretty={pretty}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return response.text, None

        return None, response.json()["error"]

    def get_gene_annotation_text(bgc_id: str, annotation_id: int, pretty=False):
        pretty = str(pretty).lower()
        response = requests.get(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/annotation/{annotation_id}?pretty={pretty}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return response.text, None

        return None, response.json()["error"]

    def update_or_create_gene_addition(bgc_id: str, data_json):
        # have to fix the strand here
        data_json["location"]["strand"] = int(data_json["location"]["strand"])

        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/to_add",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data_json,
        )

        if response.status_code == 200:
            return response.json(), None

        return None, response.json()["error"]

    def update_or_create_gene_deletion(bgc_id: str, data_json):
        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/to_delete",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data_json,
        )

        if response.status_code == 200:
            return response.json(), None

        return None, response.json()["error"]

    def update_or_create_gene_annotation(bgc_id: str, data_json):
        response = requests.post(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/annotation",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data_json,
        )

        if response.status_code == 200:
            return response.json(), None

        return None, response.json()["error"]

    def remove_gene_addition(bgc_id, addition_id: int):
        response = requests.delete(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/to_add/{addition_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return None

        return response.json()["error"]

    def remove_gene_deletion(bgc_id, deletion_id: int):
        response = requests.delete(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/to_delete/{deletion_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return None

        return response.json()["error"]

    def remove_gene_annotation(bgc_id, annotation_id: int):
        response = requests.delete(
            f"{current_app.config['API_BASE']}/entry/{bgc_id}/gene_information/annotation/{annotation_id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code == 200:
            return None

        return response.json()["error"]

    @staticmethod
    def get_or_create(bgc_id: str) -> "Entry":
        """Get an entry from the database or create one if it does not exist

        Args:
            bgc_id (str): BGC identifier

        Returns:
            Entry: new or existing database entry
        """
        if bgc_id == "new":
            entry = Entry.create(bgc_id=bgc_id)
        # entry = Entry.get(bgc_id=bgc_id)
        return entry

    # TODO: save all important data
    @staticmethod
    def submit(data: dict[str, Any]) -> str:
        """Submit a new submission to the API

        Args:
            bgc_id (str): BGC identfier
            data (dict): Minimal information to save
        """

        for compound in data["compounds"]:
            if "mass" in compound:
                compound["mass"] = float(compound["mass"])

        response = requests.post(
            f"{current_app.config['API_BASE']}/submission",
            headers={"Authorization": f"Bearer {session['token']}"},
            json=data,
        )

        if response.status_code != 200:
            return None

        return response.json()

    def check_lock(bgc_id: str, category: str):
        lock_endpoint = "/lock/check/"
        response = requests.post(
            f"{current_app.config['API_BASE']}" + lock_endpoint,
            headers={"Authorization": f"Bearer {session['token']}"},
            json={
                "accession": bgc_id,
                "category": category,
            },
        )

        return response

    def request_lock(bgc_id: str, category: str):
        lock_endpoint = "/lock/request/"
        response = requests.post(
            f"{current_app.config['API_BASE']}" + lock_endpoint,
            headers={"Authorization": f"Bearer {session['token']}"},
            json={
                "accession": bgc_id,
                "category": category,
            },
        )

        return response

    def release_lock(bgc_id: str, category: str):
        lock_endpoint = "/lock/release/"
        response = requests.post(
            f"{current_app.config['API_BASE']}" + lock_endpoint,
            headers={"Authorization": f"Bearer {session['token']}"},
            json={
                "accession": bgc_id,
                "category": category,
            },
        )

        return response

    @staticmethod
    def save_structure(bgc_id: str, data: dict[str, Any]):
        """Save structure information

        Args:
            bgc_id (str): BGC identfier
            data (dict): structure information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        for structure in data["structures"]:
            for evidence in structure["evidence"]:
                refs.update(evidence["references"])

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    @staticmethod
    def save_activity(bgc_id: str, data: dict[str, Any]):
        """Save activity information

        Args:
            bgc_id (str): BGC identfier
            data (dict): activity information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        for activity in data["activities"]:
            if assays := activity.get("assays"):
                for assay in assays:
                    refs.update(assay["references"])

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    @staticmethod
    def save_biosynth(bgc_id: str, b_class: str, data: dict[str, Any]):
        """Save biosynth information

        Args:
            bgc_id (str): BGC identfier
            b_class (str): biosynthetic class
            data (dict): biosynth information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        if b_class == "NRPS":
            for rel_type in data["release_types"]:
                if references := rel_type.get("references"):
                    refs.update(references)
        elif b_class == "Saccharide":
            for glyc in data["glycosyltransferases"]:
                refs.update(glyc["references"])
            if subcl := data.get("subclusters"):
                for sub in subcl:
                    if references := sub.get("references"):
                        refs.update(references)

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    @staticmethod
    def save_biosynth_paths(bgc_id: str, data: dict[str, Any]):
        """Save biosynthetic path information

        Args:
            bgc_id (str): BGC identfier
            data (dict): biosynthetic path information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        for path in data["paths"]:
            refs.update(path["references"])

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    @staticmethod
    def save_tailoring(bgc_id: str, data: dict[str, Any]):
        """Save tailoring information

        Args:
            bgc_id (str): BGC identfier
            data (dict): tailoring information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        if enzymes := data.get("enzymes"):
            for enzyme in enzymes:
                refs.update(enzyme["enzyme"][0]["references"])
                if reactions := enzyme.get("reactions"):
                    for reaction in reactions:
                        for evidence in reaction["reaction_smarts"][0]["evidence_sm"]:
                            refs.update(evidence["references"])

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    @staticmethod
    def save_annotation(bgc_id: str, data: dict[str, Any]):
        """Save annotation information

        Args:
            bgc_id (str): BGC identfier
            data (dict): annotation information to save
        """
        entry = Entry.get_or_create(bgc_id=bgc_id)

        refs = set()
        if annotations := data.get("annotations"):
            for annotation in annotations:
                if functions := annotation.get("functions"):
                    for function in functions:
                        if refers := function["mutation_phenotype"].get("references"):
                            refs.update(refers)
                        for evidence in function["evidence"]:
                            refs.update(evidence["references"])

        if domains := data.get("domains"):
            for domain in domains:
                if substrates := domain.get("substrates"):
                    for substrate in substrates:
                        if evidences := substrate.get("evidence"):
                            for evidence in evidences:
                                refs.update(evidences["references"])

        loaded_refs = Reference.load_missing(list(refs))
        entry.add_references(loaded_refs)

    # TODO: create batch entries from data dir
