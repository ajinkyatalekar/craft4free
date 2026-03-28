# https://craft4free.com 
This repo contains the react codebase for the web frontend for craft4free. This is hosted on AWS Amplify with CI/CD pipelines for dev and main branches.  
  
main: https://main.d2w3788e6h0dd9.amplifyapp.com/  
dev: https://dev.d2w3788e6h0dd9.amplifyapp.com/  

## Technologies
- Vite React (TypeScript)
- TailwindCSS
- ShadCN UI

## Setup Development Env
1. Clone repository.
2. Add .env with Supabase variables to the root.
3. `npm install`
4. `npm run dev`

## Development scripts
- `npm run update-types`: Update _database.types.ts_ to newer database schema.
- `openapi-typescript https://api.craft4free.com/openapi.json -o api.types.ts`: Update _api.types.ts_ to newer fastapi server types.
- `npm run preview`: Run to make sure typescript code compiles for deployment.

## Development Standards
- Craft4Free follows the standard Vite React file structure.
- _src/pages_ holds all the pages and _src/components_ holds all the components.
- Only _landing.tsx_ and _src/Auth/*_ are unprotected and all other routes require supabase authentication.
