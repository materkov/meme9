const store: any = {};
/** @ts-ignore */
window.store = store;

export function writeStore(item: any) {
    if (typeof item === "object") {
        const itemFlat: any = {};
        for (const key in item) {
            if (typeof item[key] !== "object") {
                itemFlat[key] = item[key];
            } else {
                writeStore(item[key]);
            }
        }

        store[item['id']] = itemFlat;
    }
}
