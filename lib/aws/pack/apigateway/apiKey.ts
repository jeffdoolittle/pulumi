// *** WARNING: this file was generated by the Lumi IDL Compiler (LUMIDL). ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as lumi from "@lumi/lumi";

import {RestAPI} from "./restAPI";
import {Stage} from "./stage";

export class APIKey extends lumi.Resource implements APIKeyArgs {
    public readonly name: string;
    public readonly keyName?: string;
    public description?: string;
    public enabled?: boolean;
    public stageKeys?: StageKey;

    constructor(name: string, args: APIKeyArgs) {
        super();
        if (name === undefined) {
            throw new Error("Missing required resource name");
        }
        this.name = name;
        this.keyName = args.keyName;
        this.description = args.description;
        this.enabled = args.enabled;
        this.stageKeys = args.stageKeys;
    }
}

export interface APIKeyArgs {
    readonly keyName?: string;
    description?: string;
    enabled?: boolean;
    stageKeys?: StageKey;
}

export interface StageKey {
    restAPI?: RestAPI;
    stage?: Stage;
}


