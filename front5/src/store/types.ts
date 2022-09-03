export type Composer = {}

export type Feed = {
    posts?: string[];
    nodes?: Nodes;
    route?: string;
}

export type Nodes = {
    posts?: Post[];
    users?: User[];
}

export type PostPage = {
    pagePost?: string;
    nodes?: Nodes;
}

export type UserPage = {
    pageUser?: string;
    posts?: string[];
    notFound?: boolean;
    nodes?: Nodes;
}

export type User = {
    id: string;
    name?: string;
    href?: string;
}

export type Post = {
    id: string;
    fromId?: string;
    text?: string;
    detailsURL?: string;
}

export type BrowseResult = {
    feed?: Feed;
    userPage?: UserPage;
    postPage?: PostPage;
}
