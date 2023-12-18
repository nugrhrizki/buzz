import { createMutation, createQuery } from "@tanstack/solid-query";

import { request } from "@/lib/request";

import { CreateRole, Role } from "@/models/role";

const queryRole = () => request.get<Array<Role>>("/api/v1/role/get-all");
const createQueryRole = (role: CreateRole) => request.post<Role>("/api/v1/role/create", role);
const updateQueryRole = (id: number, role: CreateRole) => request.put<Role>(`/api/v1/role/update/${id}`, role);
const deleteQueryRole = (id: number) => request.delete<Role>(`/api/v1/role/delete/${id}`);

export function useRoles() {
  return createQuery(() => ({
    queryKey: ["roles"],
    queryFn: async () => {
      const response = await queryRole();
      return response.data;
    },
  }));
}

export function createRoleMutation() {
  return createMutation(() => ({
    queryKey: ["roles"],
    mutationKey: ["roles"],
    mutationFn: async (roles: CreateRole) => {
      const response = await createQueryRole(roles);
      return response.data;
    },
  }));
}

export function updateRoleMutation() {
  return createMutation(() => ({
    queryKey: ["roles"],
    mutationKey: ["roles"],
    mutationFn: async (props: { id: number; role: CreateRole }) => {
      const response = await updateQueryRole(props.id, props.role);
      return response.data;
    },
  }));
}

export function deleteRoleMutation() {
  return createMutation(() => ({
    queryKey: ["roles"],
    mutationKey: ["roles"],
    mutationFn: async (id: number) => {
      const response = await deleteQueryRole(id);
      return response.data;
    },
  }));
}
