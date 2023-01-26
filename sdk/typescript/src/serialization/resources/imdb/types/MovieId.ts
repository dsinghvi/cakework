/**
 * This file was auto-generated by Fern from our API Definition.
 */

import * as serializers from "../../..";
import { JessieYoungApi } from "../../../..";
import * as core from "../../../../core";

export const MovieId: core.serialization.Schema<
  serializers.MovieId.Raw,
  JessieYoungApi.MovieId
> = core.serialization.string();

export declare namespace MovieId {
  type Raw = string;
}
