import Datatable from "@/components/ui/datatable";

function LogPage() {
  return (
    <div class="space-y-4 p-8 pt-6">
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight">Log</h2>
      </div>
      <Datatable.Root>
        <Datatable.Table />
      </Datatable.Root>
    </div>
  );
}
export default LogPage;
