import {Global, Page} from "../store";

export interface AppendFeed {
    type: 'feed/append'
    items: string[]
    prepend?: boolean
    nextCursor: string
}

export function appendFeed(state: Global, data: AppendFeed): Global {
    const newItem: Page = {
        items: [...data.items],
        nextCursor: data.nextCursor,
    }

    return {
        ...state,
        feed: {
            ...state.feed,
            pages: data.prepend ? [newItem, ...state.feed.pages] : [...state.feed.pages, newItem],
        }
    }
}

export interface DeleteFromFeed {
    type: 'feed/delete'
    postId: string
}

export function deleteFromFeed(state: Global, data: DeleteFromFeed): Global {
    let pages = [];
    for (let page of state.feed.pages) {
        const newPage: Page = {
            ...page,
            items: page.items.filter(item => item != data.postId)
        };
        if (newPage.items.length > 0) {
            pages.push(newPage);
        }
    }

    return {
        ...state,
        feed: {
            ...state.feed,
            pages: pages,
        }
    }
}
