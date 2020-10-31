import * as schema from "../../schema/login";
import React from "react";

interface UnathorizedPlaceholderProps {
    data: schema.ComposerRenderer;
}

export function UnathorizedPlaceholder(props: UnathorizedPlaceholderProps) {
    return <div>{props.data.unathorizedText}</div>
}