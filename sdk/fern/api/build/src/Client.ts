/**
 * This file was auto-generated by Fern from our API Definition.
 */

import { Client as ImdbClient } from "./resources/imdb/client/Client";

export declare namespace JessieYoungApiClient {
  interface Options {
    environment: string;
  }
}

export class JessieYoungApiClient {
  constructor(private readonly options: JessieYoungApiClient.Options) {}

  #imdb: ImdbClient | undefined;

  public get imdb(): ImdbClient {
    return (this.#imdb ??= new ImdbClient(this.options));
  }
}
