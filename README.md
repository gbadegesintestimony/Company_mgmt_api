Company Management API

A practical, end-to-end backend project that demonstrates how to build a secure multi-tenant company management system where an administrator can create and manage a company and its employees while employees authenticate, reset passwords, and manage their own profiles.

This project focuses on mastering database design, authentication flows, RBAC authorization, CRUD operations, migrations, and email-based OTP workflows.

Table of Contents

Overview

Project Objectives

System Architecture

Technology Stack

Features

Data Model

RBAC and Security Model

Authentication and Email Flows

API Endpoints

Project Structure

Setup and Installation

Environment Variables

Running the Project

Database Migrations

Testing the API

Future Improvements

Overview

The Company Management API is a backend service designed to simulate a real-world SaaS multi-tenant environment where multiple companies can exist within the same system while maintaining strict data isolation.

Each company has:

An Admin

Multiple Employees

The admin manages company resources and employees, while employees manage their own profiles and authentication.

The system includes:

Secure authentication

Role-based authorization

Email verification

OTP-based password reset

Session management with refresh tokens

Audit logging

Project Objectives
Primary Goal

Build a secure multi-tenant backend system where administrators can manage company operations and employees securely.

Learning Outcomes

This project helps internalize the following backend engineering concepts:

Database Fundamentals

Schema design

PostgreSQL relational modeling

Migrations

Indexes

Joins

Transactions

Authentication

JWT access tokens

Refresh token rotation

Session persistence

Password hashing

Authorization

Role Based Access Control (RBAC)

Company-scoped authorization

CRUD Operations

Create, read, update, and delete across:

Companies

Employees

Profiles

Email and OTP Systems

Render

OTP verification

Password reset workflows

Rate limiting

Security

Password hashing (bcrypt)

Input validation

Rate limiting

Least privilege access

Observability

Structured logging

Request tracing

Error handling

System Architecture
Client (Postman / Web App)
│
▼
Go HTTP API (REST)
│
▼
Service Layer
(Authentication, RBAC, Business Logic)
│
▼
Repository Layer
(Database Queries)
│
▼
PostgreSQL Database
Architecture Components

API Layer

Handles HTTP requests and responses using Go's HTTP server.

Authentication Layer

Handles:

JWT access tokens

Refresh token rotation

Session validation

RBAC Middleware

Restricts access based on role:

Admin

Employee

Database Layer

PostgreSQL with migrations for schema evolution.

Email Layer

Render integration used for:

Email verification

Password reset OTPs

Configuration

Environment variables store secrets such as:

JWT keys

Render credentials

Database connection strings

Observability

Logging system for request tracking and error visibility.

Technology Stack
Component Technology
Language Go
HTTP Framework net/http + router
Database PostgreSQL
Migrations golang-migrate
Authentication JWT
Password Hashing bcrypt / argon2id
Email Render
Containerization Docker
Testing Postman
Features
Admin Features

Create company

Manage employees

Update employee roles

Deactivate / reactivate users

View audit logs

Manage company profile

Employee Features

Login

Update profile

Change password

Reset password via OTP

View company information

Security Features

JWT authentication

Refresh token sessions

Password hashing

OTP verification

Role-based access control

Company-scoped authorization

Data Model
Companies
Field Description
id Company ID
name Company name
domain Company email domain
status Active / suspended
created_at Creation timestamp
Users
Field Description
id User ID
company_id Foreign key to company
role admin or employee
email User email
password_hash Hashed password
is_active Account status
email_verified_at Email verification timestamp
Profiles

Stores employee profile information.

Field Description
user_id FK to user
first_name First name
last_name Last name
phone Phone number
job_title Job title
department Department
Password Resets

Stores OTPs for password reset.

Field Description
user_id User reference
otp_code Reset code
expires_at Expiration time
Email Verifications

Stores email verification OTPs.

Sessions

Stores refresh token sessions.

Field Description
refresh_token_hash Stored hashed token
expires_at Session expiry
Audit Logs

Tracks system events.

Example:

Employee created

Employee deleted

Password reset

Role updated

RBAC and Security Model
Roles

Admin

Manage company

Create employees

Update employees

View audit logs

Employee

Login

Update own profile

Reset password

View company details

Authorization Rules

Every request is company scoped

Employees cannot access admin endpoints

Inactive users cannot authenticate

Unverified users cannot access privileged endpoints

Authentication and Email Flows
Admin Onboarding

Admin registers company

Email verification OTP sent

Admin confirms OTP

Admin session created

Employee Lifecycle

Admin creates employee

Employee receives invite email

Employee verifies email via OTP

Employee logs in

Employee updates profile

Password Reset

User requests reset

OTP sent via email

User confirms OTP

User sets new password

API Endpoints
Authentication

POST /v1/auth/admin/register

Create company and admin account.

POST /v1/auth/admin/login

Admin login.

POST /v1/auth/employee/login

Employee login.

POST /v1/auth/refresh

Refresh JWT token.

POST /v1/auth/logout

Invalidate refresh token.

Email Verification

POST /v1/verification/email/request

Request verification OTP.

POST /v1/verification/email/confirm

Verify email with OTP.

Password Reset

POST /v1/password/forgot

Request password reset OTP.

POST /v1/password/reset

Reset password using OTP.

Company Management

GET /v1/companies/:companyId

Get company details.

PATCH /v1/companies/:companyId

Update company details.

Employee Administration

POST /v1/companies/:companyId/employees

Create employee.

GET /v1/companies/:companyId/employees

List employees.

GET /v1/companies/:companyId/employees/:userId

View employee.

PATCH /v1/companies/:companyId/employees/:userId

Update employee.

DELETE /v1/companies/:companyId/employees/:userId

Delete employee.

Employee Self Service

GET /v1/me

Get current user.

PATCH /v1/me/profile

Update own profile.

PATCH /v1/me/password

Change password.

Project Structure
company_mgmt_api

cmd/
api/

internal/
handlers/
services/
repositories/
middleware/

database/
migrations/

pkg/
auth/
email/
utils/

docker-compose.yml
.env
Setup and Installation
Clone repository
git clone https://github.com/yourname/company-mgmt-api
cd company-mgmt-api
Start services
docker compose up --build

This will start:

API server

PostgreSQL database

Environment Variables

Example .env

PORT=8080

DB_HOST=postgres
DB_PORT=5432
DB_USER=company_user
DB_PASSWORD=company_pass
DB_NAME=company_db

JWT_SECRET=supersecret

RESEND_API_KEY
EMAIL_FROM

Database Migrations

This project uses golang-migrate.

Run migrations:

migrate up

Rollback:

migrate down

Migration files are located in:

database/migrations
Testing the API

You can test endpoints using:

Postman

Curl

REST clients

Example:

POST http://localhost:8080/v1/auth/admin/register
Future Improvements

Possible enhancements:

API documentation with Swagger

Rate limiting middleware

Background job queue

Email templates

Metrics and monitoring

CI/CD pipeline

Unit and integration testing

OAuth authentication
