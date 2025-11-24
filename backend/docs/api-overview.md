# SportLink Core Service - API Overview

## What is SportLink Core Service?

SportLink Core Service is a RESTful API designed to facilitate the organization and management of sports activities. The service enables users to connect with others interested in playing specific sports and helps coordinate matches by managing players and teams.

## Business Domain

This API operates in the **sports management and social networking** domain, focusing on bringing together people who want to participate in sports activities. It currently supports multiple sports including Football, Paddle, and Tennis.

## Intended Users

The API is designed for:
- **Mobile and Web Applications**: Frontend applications that need to manage sports teams and players
- **Sports Facility Management Systems**: Platforms that coordinate sports activities and venues
- **Community Sports Organizers**: Tools that help organize amateur and recreational sports events
- **Developers**: Building sports-related social and organizational features

## Main Capabilities

The SportLink Core Service provides the following core capabilities:

1. **Team Management**: Create and retrieve sports teams with specific characteristics such as sport type, category level, and team members
2. **Player Management**: Handle individual player profiles with sport preferences and skill levels
3. **Multi-Sport Support**: Support for various sports including Paddle, Football, and Tennis
4. **Skill-Based Categorization**: Teams and players can be categorized by skill level (Unranked, L1-L7)
5. **Member Validation**: Ensures team members exist in the system before creating teams

## Problem It Solves

SportLink Core Service addresses several key challenges:

- **Finding Compatible Players**: Helps match players with similar skill levels and sport interests
- **Team Organization**: Simplifies the process of creating and managing sports teams
- **Skill Level Management**: Maintains structured skill categories to ensure balanced competition
- **Data Consistency**: Validates team composition to ensure all members are properly registered
- **Sport-Specific Organization**: Allows filtering and organization by specific sports

## Technical Approach

The API follows **Hexagonal Architecture** principles, ensuring:
- Clean separation between business logic and infrastructure concerns
- Easy testability and maintainability
- Flexibility to adapt to different data sources and interfaces
- Domain-driven design with clear bounded contexts

## Supported Sports

Currently, the API supports:
- **Paddle** (Pádel)
- **Football** (Soccer)
- **Tennis**

## Skill Categories

Teams and players are classified using a standardized category system:
- **Unranked**: No skill level assigned
- **L1 to L7**: Progressive skill levels from beginner (L1) to advanced (L7)

This categorization helps in creating balanced teams and organizing appropriate matches.

