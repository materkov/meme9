import * as schema from "../../schema/api";
import React from "react";

interface UnathorizedPlaceholderProps {
    data: schema.ComposerRenderer;
}

export function UnathorizedPlaceholder(props: UnathorizedPlaceholderProps) {
    return <div>{props.data.unathorizedText}</div>
}