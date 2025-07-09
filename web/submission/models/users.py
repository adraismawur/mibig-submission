from __future__ import annotations
from dataclasses import asdict, dataclass
from enum import Enum
from flask import current_app, json, session
from flask_login import UserMixin
import requests


from submission.extensions import db


class Role(Enum):
    SUBMITTER = 0
    REVIEWER = 1
    ADMIN = 2


@dataclass
class UserRole:
    role: Role

    def asdict(self):
        return asdict(self)


@dataclass
class UserInfo:
    alias: str
    name: str
    call_name: str
    orc_id: str
    organisation_1: str
    organisation_2: str
    organisation_3: str
    public: bool

    @staticmethod
    def generate_alias(length: int = 15) -> str:
        import base64
        from random import randbytes

        token_bytes = randbytes(length)

        return base64.b32encode(token_bytes).decode("utf-8")

    @staticmethod
    def guess_call_name(name: str) -> str:
        return name.split(" ")[0]

    @staticmethod
    def from_json(info_dict):
        name = info_dict["name"]
        alias = info_dict["alias"]
        call_name = info_dict["call_name"]
        orc_id = info_dict["orc_id"]
        organisation = info_dict["organization1"]
        organisation_2 = info_dict["organization2"]
        organisation_3 = info_dict["organization3"]
        public = info_dict["public"]

        return UserInfo(
            alias,
            name,
            call_name,
            orc_id,
            organisation,
            organisation_2,
            organisation_3,
            public,
        )

    def asdict(self):
        return asdict(self)


class User(UserMixin):
    id: int
    email: str
    active: bool

    roles: list[UserRole]

    info: UserInfo

    def __init__(
        self, id: int, email: str, active: bool, roles: list[UserRole], info: UserInfo
    ):
        self.id = id
        self.email = email
        self.active = active
        self.roles = roles
        self.info = info

    @staticmethod
    def from_json(json_dict):

        # assemble roles
        roles = []

        for role_dict in json_dict["roles"]:
            roles.append(UserRole(role_dict["role"]))

        info_dict = json_dict["info"]
        info = UserInfo.from_json(info_dict)

        user = User(
            json_dict["id"],
            json_dict["email"],
            json_dict["active"],
            roles,
            info,
        )

        return user

    @staticmethod
    def get_user(id: int) -> User:
        response = requests.get(
            f"{current_app.config['API_BASE']}/user/{id}",
            headers={"Authorization": f"Bearer {session['token']}"},
        )

        if response.status_code != 200:
            return None

        user_data = response.json()

        return User.from_json(user_data)

    def has_role(self, role: Role):
        for userRole in self.roles:
            if userRole.role == role:
                return True

        return False

    def is_admin(self):
        return self.has_role(Role.ADMIN)

    def to_json(self):
        return json.dumps(
            {
                "email": self.email,
                "active": self.active,
                "roles": [r.asdict() for r in self.roles],
                "info": self.info.asdict(),
            }
        )
