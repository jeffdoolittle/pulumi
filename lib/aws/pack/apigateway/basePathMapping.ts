// *** WARNING: this file was generated by the Lumi IDL Compiler (LUMIDL). ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as lumi from "@lumi/lumi";

import {RestAPI} from "./restAPI";
import {Stage} from "./stage";

export class BasePathMapping extends lumi.Resource implements BasePathMappingArgs {
    public readonly name: string;
    public domainName: string;
    public restAPI: RestAPI;
    public basePath?: string;
    public stage?: Stage;

    constructor(name: string, args: BasePathMappingArgs) {
        super();
        if (name === undefined) {
            throw new Error("Missing required resource name");
        }
        this.name = name;
        if (args.domainName === undefined) {
            throw new Error("Missing required argument 'domainName'");
        }
        this.domainName = args.domainName;
        if (args.restAPI === undefined) {
            throw new Error("Missing required argument 'restAPI'");
        }
        this.restAPI = args.restAPI;
        this.basePath = args.basePath;
        this.stage = args.stage;
    }
}

export interface BasePathMappingArgs {
    domainName: string;
    restAPI: RestAPI;
    basePath?: string;
    stage?: Stage;
}


