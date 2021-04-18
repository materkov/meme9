declare module "*.module.css";

interface Window {
    modules: { [name: string]: any };
}
