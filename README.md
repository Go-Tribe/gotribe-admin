[English](README.md) | [中文](README_CN.md)

---
<h1 align="center">gotribe-admin</h1>

<div align="center">
`gotribe-admin` is a small CMS solution developed with Go and Vue, featuring a rich set of themes, ready-to-use out of the box, and an enterprise-level architecture. It is suitable for individuals, teams, and small to medium-sized enterprises.
<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/go-tribe/gotribe-admin" alt="Go version"/>
<img src="https://img.shields.io/badge/Gin-1.9.1-brightgreen" alt="Gin version"/>
<img src="https://img.shields.io/badge/Gorm-1.25.8-brightgreen" alt="Gorm version"/>
<img src="https://img.shields.io/github/license/go-tribe/gotribe-admin" alt="License"/>
</p>
</div>

### Core Advantages

- **Performance**: Leveraging Golang's efficient concurrency capabilities, GoTribe can easily handle the demands of high-traffic websites.
- **Ease of Use**: With a clean and intuitive user interface and documentation, even beginners can get started quickly.
- **Highly Customizable**: Offers a wealth of APIs and plugin support to meet personalized website building needs.
- **Community Support**: An active open-source community provides ongoing updates and technical support.
- **Security and Stability**: Adheres to best security practices to ensure the safety and stable operation of website data.

### Applicable Scenarios

GoTribe provides robust support and flexible customization options for everything from personal blogs and team projects to enterprise websites.

### Demo

![Login](https://github.com/Go-Tribe/gotribe-admin/blob/main/docs/images/login.png)
![Dashboard](https://github.com/Go-Tribe/gotribe-admin/blob/main/docs/images/index.png)
![System Management](https://github.com/Go-Tribe/gotribe-admin/blob/main/docs/images/system.png)
![Log Management](https://github.com/Go-Tribe/gotribe-admin/blob/main/docs/images/log.png)
![Project Management](https://github.com/Go-Tribe/gotribe-admin/blob/main/docs/images/project.png)
![Content Management](https://github.com/Go-Tribe/gotribe-admin/blob/main/docs/images/content.png)

### Project Description

The project adopts a front-end and back-end separation architecture, consisting of three parts: the management API, the business API, and the management backend UI. The business front-end UI can be developed according to your needs or using our templates.

#### Projects

| Project              | Description  | Address                             |
|----------------------|--------------|--------------------------------------|
| **gotribe-admin**    | Management API| [Link](https://github.com/go-tribe/gotribe-admin.git) |
| **gotribe**          | Business API | [Link](https://github.com/go-tribe/gotribe.git)     |
| **gotribe-admin-vue**| Management UI | [Link](https://github.com/go-tribe/gotribe-admin-vue.git) |

#### Business Themes

| Theme             | Description    | Address                                      |
|-------------------|----------------|----------------------------------------------|
| **gotribe-blog**  | A simple blog theme | [Link](https://github.com/go-tribe/gotribe-blog.git) |

#### Relationship Diagram

```mermaid
graph LR
    A[Go-Tribe Project] --> B(gotribe-admin Management Backend)
    A --> C(gotribe Business API)
    A --> E(gotribe-blog Blog Theme)

    B --> F[Database]
    C --> F

    E --> G[Business Frontend UI]
    G -->|User Customized| H[Business Themes]
    H --> I[gotribe-blog Blog Theme]
```

The diagram above clearly illustrates the structure of the Go-Tribe project and the interactions between its components:

- **Go-Tribe** is the name of the entire system framework, which includes multiple modules, each responsible for different functions.
- **gotribe-admin Management Backend**: This is the core management module of the system, used for handling backend management tasks. For security reasons, it is usually deployed on an internal network and accessed via VPN. To simplify the deployment process, we have integrated the gotribe-admin-vue Management UI with the Management API for one-click deployment.
- **gotribe Business API**: This module is responsible for handling business logic, with a particular focus on search engine optimization (SEO) and development efficiency. It is completely decoupled from the business frontend UI and supports Kubernetes deployment and horizontal scaling to accommodate the needs of businesses of different sizes.
- **gotribe-blog Blog Theme**: Provides a pre-built blog theme as an example of a business theme, demonstrating how to quickly build specific business scenarios using the Go-Tribe framework.
- **Database**: Serves as the data storage center of the system, responsible for saving all necessary data.
- **Business Frontend UI**: Users can develop customized frontend interfaces according to their specific needs, using templates provided by Go-Tribe.

The entire system is designed with a front-end and back-end separation architecture, which not only improves the system's flexibility but also allows each component to be developed and maintained independently, thereby enhancing the system's scalability and maintainability.

### Quick Start

> Prerequisites: `go1.20+` `node 18+`.

1. Download the project

```bash
git clone --recursive https://github.com/go-tribe/gotribe-admin.git
```
2. Modify the configuration file, parameter details are in: `config.tmp.yml`.

```bash
cp config.tmp.yml config.yml
```
3. Run for development:

```bash
make run
```
4. Compile and package:

```bash
make
```
After compilation, the executable files will be in the `_output` directory.

The database will be migrated automatically upon startup, and the default username for the management backend is: admin, with the password: 123456.

### TODO

- Add payment configuration
- Add product management

### Online Application

[Mafan](https://www.dengmengmian.com)

### License

[MIT](https://choosealicense.com/licenses/mit/)

---

Please note that the image and link references in the documentation are placeholders and should be replaced with the actual URLs when the document is compiled.
