# go-filerX

This project was born out of frustration when I discovered that my favorite PDF hosting site no longer allowed free downloads and many of the PDFs I relied on had been deleted. Perhaps they were reported by competitors.

So, I decided to create my own PDF uploader site where everyone can freely download and upload documents of all sorts.

## Rebuilt core

The application has been rebuilt to be runnable out of the box without a database dependency:

- Upload files through a web form (`/`).
- List and download uploaded files (`/listpdf`).
- Delete uploads from the list page.
- Register and login users with bcrypt-hashed passwords kept in memory (with a local JSON snapshot for backup).
- Basic validation for usernames/passwords and upload type/size.

## Run locally

```bash
go run .
```

Then open: `http://localhost:8080`

Uploads are stored in the `uploads/` directory.
