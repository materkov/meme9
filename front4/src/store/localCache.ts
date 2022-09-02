import {getByType} from "./store";

function saveToCache() {
    const query = getByType("Query");
    if (!query || query.type != "Query") {
        return;
    }

    localStorage.setItem("viewerId", query.viewer || "");
}
