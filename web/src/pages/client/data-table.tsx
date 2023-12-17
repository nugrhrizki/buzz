import {
  ColumnDef,
  ColumnFiltersState,
  SortingState,
  VisibilityState,
  createSolidTable,
  getCoreRowModel,
  getExpandedRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
} from "@tanstack/solid-table";
import { Match, Switch, createSignal } from "solid-js";

import { Client } from "@/models/client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import DataTable from "@/components/ui/data-table";
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table";

type ColumnDefiniton = ColumnDef<Client>;

export const columns: ColumnDefiniton[] = [
  DataTable.RowExpand(),
  {
    accessorKey: "name",
    header: (header) => <DataTable.ColumnHeader column={header.column} title="Nama" />,
    enableHiding: false,
  },
  {
    accessorKey: "token",
    header: (header) => <DataTable.ColumnHeader column={header.column} title="Token" />,
  },
  {
    accessorKey: "connected",
    header: (header) => <DataTable.ColumnHeader column={header.column} title="Status" />,
    cell: (cell) => (
      <Switch fallback={<Badge variant="secondary">-</Badge>}>
        <Match when={cell.row.original.connected === 1}>
          <Badge variant="default">connected</Badge>
        </Match>
        <Match when={cell.row.original.connected === 0}>
          <Badge variant="destructive">disconnected</Badge>
        </Match>
      </Switch>
    ),
  },
];

export function createClientTable(data: Client[], columns: ColumnDefiniton[]) {
  const [sorting, setSorting] = createSignal<SortingState>([]);
  const [columnFilters, setColumnFilters] = createSignal<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = createSignal<VisibilityState>({});
  const [rowSelection, setRowSelection] = createSignal({});

  return createSolidTable({
    data,
    columns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getRowCanExpand: () => true,
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      get sorting() {
        return sorting();
      },
      get columnFilters() {
        return columnFilters();
      },
      get columnVisibility() {
        return columnVisibility();
      },
      get rowSelection() {
        return rowSelection();
      },
    },
  });
}

interface TableDetailProps {
  client: Client;
}

export function TableDetail(props: TableDetailProps) {
  return (
    <div class="p-4">
      <Card>
        <CardContent class="p-2 pt-2">
          <Table>
            <TableBody>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">ID</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.client.id || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">JID</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.client.jid || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Name</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.client.name || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Token</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>{props.client.token || "-"}</TableCell>
              </TableRow>
              <TableRow>
                <TableCell class="w-[1%] whitespace-nowrap">Status</TableCell>
                <TableCell class="w-[1%] whitespace-nowrap">:</TableCell>
                <TableCell>
                  <Switch fallback={<Badge variant="secondary">-</Badge>}>
                    <Match when={props.client.connected === 1}>
                      <Badge variant="default">connected</Badge>
                    </Match>
                    <Match when={props.client.connected === 0}>
                      <Badge variant="destructive">disconnected</Badge>
                    </Match>
                  </Switch>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
        <CardFooter>
          <div class="flex justify-end w-full gap-x-4">
            <Button size="sm">Open</Button>
            <Button variant="outline" size="sm">
              Edit
            </Button>
            <Button variant="destructive" size="sm">
              Delete
            </Button>
          </div>
        </CardFooter>
      </Card>
    </div>
  );
}
