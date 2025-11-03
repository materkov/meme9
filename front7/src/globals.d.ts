declare module '*.module.css' {
  const classes: { [key: string]: string };
  export default classes;
}

declare global {
  interface Window {
    API_BASE_URL: string;
  }
}

export {};

