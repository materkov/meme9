import {QueryParams} from "../types";
import {api} from "../api";
import {storeOnChanged} from "./store";

export function getVkAuthUrl() {
    const urlQuery: QueryParams = {
        vkAuthUrl: {}
    };
    api(urlQuery).then(result => {
        // @ts-ignore
        window.store["fake-id-vk-auth-url"] = {
            type: "VkAuthURL",
            id: "fake-id-vk-auth-url",
            url: result.vkAuthUrl,
        }
        storeOnChanged();
    })
}

export function getViewer() {
    const q: QueryParams = {
        viewer: {
            inner: {
                name: {}
            }
        }
    }
    api(q).then(result => {
        // @ts-ignore
        window.store[result.viewer.id] = {...window.store[result.viewer.id], ...result.viewer};
        storeOnChanged();
    })
}
