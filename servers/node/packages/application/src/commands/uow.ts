import type { WriteUnitOfWork, WriteUnitOfWorkFactory } from "@library/domain";

export async function inWriteUnitOfWork<T>(
  factory: WriteUnitOfWorkFactory,
  work: (uow: WriteUnitOfWork) => Promise<T>
): Promise<T> {
  const uow = await factory.build();

  try {
    const result = await work(uow);
    await uow.commit();
    return result;
  } catch (error) {
    await rollbackQuietly(uow);
    throw error;
  }
}

async function rollbackQuietly(uow: WriteUnitOfWork): Promise<void> {
  try {
    await uow.rollback();
  } catch {
    // Preserve the original command failure. Rollback failures are infrastructure noise here.
  }
}
