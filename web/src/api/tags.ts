import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { env } from "../env.ts";

export type Tag = {
  id: string;
  name: string;
  createdAt: Date;
  updatedAt: Date;
};

type TagResponse = {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
};

function isTagResponse(arg: unknown): arg is TagResponse {
  const t = arg as TagResponse;
  return (
    typeof t?.id === "string" &&
    typeof t?.name === "string" &&
    typeof t?.created_at === "string" &&
    typeof t?.updated_at === "string"
  );
}

function toTag(r: TagResponse): Tag {
  return {
    id: r.id,
    name: r.name,
    createdAt: new Date(r.created_at),
    updatedAt: new Date(r.updated_at),
  };
}

type TagsResponse = {
  tags: TagResponse[];
};

function isTagsResponse(arg: unknown): arg is TagsResponse {
  const ts = arg as TagsResponse;
  return Array.isArray(ts?.tags) && ts?.tags.every(isTagResponse);
}

export function useTag(tagID: string) {
  return useQuery({
    queryKey: ["tags", tagID],
    queryFn: async () => {
      const response = await fetch(`${env.apiBaseURL}/tags/${tagID}`, {
        method: "GET",
        headers: { Authorization: `Bearer ${env.apiToken}` },
      });
      if (!response.ok) {
        throw new Error(`HTTP error status: ${response.status}`);
      }

      const data: unknown = await response.json();
      if (isTagResponse(data)) {
        return toTag(data);
      }
      throw new Error("invalid response body");
    },
  });
}

export function useTags() {
  return useQuery({
    queryKey: ["tags"],
    queryFn: async () => {
      const response = await fetch(`${env.apiBaseURL}/tags?limit=10&offset=0`, {
        method: "GET",
        headers: { Authorization: `Bearer ${env.apiToken}` },
      });
      if (!response.ok) {
        throw new Error(`HTTP error status: ${response.status}`);
      }

      const data: unknown = await response.json();
      if (isTagsResponse(data)) {
        return data.tags.map(toTag);
      }
      throw new Error("invalid response body");
    },
  });
}

type CreateTagRequest = {
  name: string;
};

export function useCreateTag() {
  const c = useQueryClient();
  return useMutation({
    mutationFn: async (req: CreateTagRequest) => {
      const response = await fetch(`${env.apiBaseURL}/tags`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${env.apiToken}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name: req.name }),
      });
      if (!response.ok) {
        throw new Error(`HTTP error status: ${response.status}`);
      }
      return await response.json();
    },
    onSuccess: () => {
      void c.invalidateQueries({ queryKey: ["tags"] });
    },
  });
}

type UpdateTagRequest = {
  tagID: string;
  name: string;
};

export function useUpdateTag() {
  const c = useQueryClient();
  return useMutation({
    mutationFn: async (req: UpdateTagRequest) => {
      const response = await fetch(`${env.apiBaseURL}/tags/${req.tagID}`, {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${env.apiToken}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name: req.name }),
      });
      if (!response.ok) {
        throw new Error(`HTTP error status: ${response.status}`);
      }
    },
    onSuccess: () => {
      void c.invalidateQueries({ queryKey: ["tags"] });
    },
  });
}

export function useDeleteTag() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: async (tagID: string) => {
      const response = await fetch(`${env.apiBaseURL}/tags/${tagID}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${env.apiToken}` },
      });
      if (!response.ok) {
        throw new Error(`HTTP error status: ${response.status}`);
      }
    },
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: ["tags"] });
    },
  });
}
