// *** WARNING: this file was generated by the Lumi IDL Compiler (LUMIDL). ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as lumi from "@lumi/lumi";

import {VPC} from "./vpc";

export class RouteTable extends lumi.Resource implements RouteTableArgs {
    public readonly name: string;
    public readonly vpc: VPC;

    constructor(name: string, args: RouteTableArgs) {
        super();
        if (name === undefined) {
            throw new Error("Missing required resource name");
        }
        this.name = name;
        if (args.vpc === undefined) {
            throw new Error("Missing required argument 'vpc'");
        }
        this.vpc = args.vpc;
    }
}

export interface RouteTableArgs {
    readonly vpc: VPC;
}


