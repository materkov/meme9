interface Window {
    modules: { [renderer: string]: object };
    rootLoaded: boolean;
    InitRequest: any;
    InitData: any;
    InitJsBundles: string[];
    apiKey: string;
    apiCache: { [key: string]: any };
}
