"""
Data models for the ArasAuth Python SDK
"""

from dataclasses import dataclass
from typing import Optional, List, Dict, Any


@dataclass
class User:
    """Represents a user in the system"""
    id: str
    email: str
    first_name: str
    last_name: str
    status: str
    email_verified: bool
    is_deleted: bool = False
    created_at: str = ''
    updated_at: str = ''

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'User':
        """Create a User instance from a dictionary"""
        return cls(
            id=data.get('id', ''),
            email=data.get('email', ''),
            first_name=data.get('first_name', ''),
            last_name=data.get('last_name', ''),
            status=data.get('status', ''),
            email_verified=data.get('email_verified', False),
            is_deleted=data.get('is_deleted', False),
            created_at=data.get('created_at', ''),
            updated_at=data.get('updated_at', '')
        )


@dataclass
class Group:
    """Represents a group in the system"""
    id: str
    name: str
    description: str
    is_active: bool = True
    is_deleted: bool = False
    created_at: str = ''
    updated_at: str = ''

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'Group':
        """Create a Group instance from a dictionary"""
        return cls(
            id=data.get('id', ''),
            name=data.get('name', ''),
            description=data.get('description', ''),
            is_active=data.get('is_active', True),
            is_deleted=data.get('is_deleted', False),
            created_at=data.get('created_at', ''),
            updated_at=data.get('updated_at', '')
        )


@dataclass
class Role:
    """Represents a role in the system"""
    id: str
    name: str
    description: str
    is_active: bool = True
    is_deleted: bool = False
    created_at: str = ''
    updated_at: str = ''

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'Role':
        """Create a Role instance from a dictionary"""
        return cls(
            id=data.get('id', ''),
            name=data.get('name', ''),
            description=data.get('description', ''),
            is_active=data.get('is_active', True),
            is_deleted=data.get('is_deleted', False),
            created_at=data.get('created_at', ''),
            updated_at=data.get('updated_at', '')
        )


@dataclass
class Permission:
    """Represents a permission in the system"""
    id: str
    resource: str
    action: str
    description: str
    is_active: bool = True
    is_deleted: bool = False
    created_at: str = ''
    updated_at: str = ''

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'Permission':
        """Create a Permission instance from a dictionary"""
        return cls(
            id=data.get('id', ''),
            resource=data.get('resource', ''),
            action=data.get('action', ''),
            description=data.get('description', ''),
            is_active=data.get('is_active', True),
            is_deleted=data.get('is_deleted', False),
            created_at=data.get('created_at', ''),
            updated_at=data.get('updated_at', '')
        )


@dataclass
class AuthResponse:
    """Represents the response from authentication endpoints"""
    access_token: str
    refresh_token: str
    expires_in: int
    token_type: str
    user: User

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'AuthResponse':
        """Create an AuthResponse instance from a dictionary"""
        user_data = data.get('user', {})
        user = User.from_dict(user_data) if user_data else None
        
        return cls(
            access_token=data.get('access_token', ''),
            refresh_token=data.get('refresh_token', ''),
            expires_in=data.get('expires_in', 0),
            token_type=data.get('token_type', 'Bearer'),
            user=user
        )


@dataclass
class ListResponse:
    """Represents a paginated list response"""
    items: List[Any]
    total: int
    page: int
    limit: int

    @classmethod
    def from_dict(cls, data: Dict[str, Any], item_class) -> 'ListResponse':
        """Create a ListResponse instance from a dictionary"""
        items_data = data.get('items', [])
        items = [item_class.from_dict(item) for item in items_data]
        
        return cls(
            items=items,
            total=data.get('total', 0),
            page=data.get('page', 1),
            limit=data.get('limit', 20)
        )


