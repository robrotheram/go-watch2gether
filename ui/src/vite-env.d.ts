interface ImportMetaEnv {
  readonly VITE_APP_VERSION: string
}

declare module '*.svg' {
    const content: string;
    export default content;
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
