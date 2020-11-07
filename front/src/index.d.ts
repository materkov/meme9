interface Window {
    modules: { [renderer: string]: object };
    rootLoaded: boolean;
    InitApiMethod: string;
    InitApiRequest: any;
    InitApiResponse: any;
    InitJsBundles: string[];
    InitRootComponent: string;
    CSRFToken: string;
    apiKey: string;
    apiCache: { [key: string]: any };
}
