import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { env } from "../env.ts";

export type Project = {
  id: string;
  name: string;
  color: string;
};

function isProject(arg: unknown): arg is Project {
  const p = arg as Project;
  return typeof p?.id === "string" && typeof p?.name === "string" && typeof p?.color === "string";
}

type Projects = {
  projects: Project[];
};

function isProjects(arg: unknown): arg is Projects {
  const ps = arg as Projects;
  return Array.isArray(ps?.projects) && ps?.projects.every(isProject);
}

export function useProject(projectID: string) {
  return useQuery({
    queryKey: ["projects", projectID],
    queryFn: async () => {
      const response = await fetch(`${env.apiBaseURL}/projects/${projectID}`, {
        method: "GET",
        headers: { Authorization: `Bearer ${env.apiToken}` },
      });
      if (!response.ok) {
        throw new Error(`HTTP error status: ${response.status}`);
      }

      const data: unknown = await response.json();
      if (isProject(data)) {
        return data;
      }
      throw new Error("invalid response body");
    },
  });
}

export function useProjects() {
  return useQuery({
    queryKey: ["projects"],
    queryFn: async () => {
      const response = await fetch(`${env.apiBaseURL}/projects?limit=10&offset=0`, {
        method: "GET",
        headers: { Authorization: `Bearer ${env.apiToken}` },
      });
      if (!response.ok) {
        throw new Error(`HTTP error status: ${response.status}`);
      }

      const data: unknown = await response.json();
      if (isProjects(data)) {
        return data.projects;
      }
      throw new Error("invalid response body");
    },
  });
}

type CreateProjectRequest = {
  name: string;
  color: string;
};

export function useCreateProject() {
  const c = useQueryClient();
  return useMutation({
    mutationFn: async (req: CreateProjectRequest) => {
      const response = await fetch(`${env.apiBaseURL}/projects`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${env.apiToken}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name: req.name, color: req.color }),
      });
      if (!response.ok) {
        throw new Error(`HTTP error status: ${response.status}`);
      }
      return await response.json();
    },
    onSuccess: () => {
      void c.invalidateQueries({ queryKey: ["projects"] });
    },
  });
}

type UpdateProjectRequest = {
  projectID: string;
  name: string;
  color: string;
};

export function useUpdateProject() {
  const c = useQueryClient();
  return useMutation({
    mutationFn: async (req: UpdateProjectRequest) => {
      const response = await fetch(`${env.apiBaseURL}/projects/${req.projectID}`, {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${env.apiToken}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name: req.name, color: req.color }),
      });
      if (!response.ok) {
        throw new Error(`HTTP error status: ${response.status}`);
      }
    },
    onSuccess: () => {
      void c.invalidateQueries({ queryKey: ["projects"] });
    },
  });
}

export function useDeleteProject() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: async (projectID: string) => {
      const response = await fetch(`${env.apiBaseURL}/projects/${projectID}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${env.apiToken}` },
      });
      if (!response.ok) {
        throw new Error(`HTTP error status: ${response.status}`);
      }
    },
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: ["projects"] });
    },
  });
}
