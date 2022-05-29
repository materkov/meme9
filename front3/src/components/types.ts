export type RootRenderer = {
    label: string;
    command: Action;
    header: HeaderRenderer;
    composer: ComposerRenderer;
}

export type HeaderRenderer = {
    userName: string;
}

export type Navigate = {
    url: string;
}

export type WriteLog = {
    message: string;
}

export type CreatePost = {
    text: string;
}

export type Action = {
    navigate?: Navigate;
    writeLog?: WriteLog;
    createPost?: CreatePost;
}

export type LoginRenderer = {
    authURL: string;
}

export type ComposerRenderer = {
    placeholder: string;
};

export type CreatePostReq = {
    text: string;
}

export type CreatePostResp = {
    id: string;
}
